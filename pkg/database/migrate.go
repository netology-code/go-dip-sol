package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Migrate применяет SQL-миграции из папки migrations
func Migrate(db *sql.DB) error {
	migrationsDir := "migrations"
	log.Printf("Looking for migrations in: %s", migrationsDir)

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Фильтруем и сортируем миграции
	var migrationFiles []string
	migrationPattern := regexp.MustCompile(`^\d+_.+\.sql$`)

	for _, file := range files {
		if !file.IsDir() && migrationPattern.MatchString(file.Name()) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Сортируем миграции по номеру
	sort.Slice(migrationFiles, func(i, j int) bool {
		return extractMigrationNumber(migrationFiles[i]) < extractMigrationNumber(migrationFiles[j])
	})

	// Проверяем, нужно ли применять миграции
	for _, file := range migrationFiles {
		log.Printf("Processing migration: %s", file)

		// Проверяем, не была ли уже применена эта миграция
		if isMigrationApplied(db, file) {
			log.Printf("Migration %s already applied, skipping", file)
			continue
		}

		// Читаем содержимое миграции
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		// Разбиваем SQL на отдельные запросы
		queries := splitSQLQueries(string(content))

		// Выполняем запросы в транзакции
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("failed to execute migration query in %s: %w", file, err)
			}
		}

		// Отмечаем миграцию как выполненную
		_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", file)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		log.Printf("Successfully applied migration: %s", file)
	}

	return nil
}

// splitSQLQueries разбивает SQL-скрипт на отдельные запросы
func splitSQLQueries(content string) []string {
	var queries []string
	var currentQuery strings.Builder
	inString := false
	stringChar := byte(0)
	inComment := false
	commentType := "" // "--" или "/*"

	for i := 0; i < len(content); i++ {
		char := content[i]

		// Обработка комментариев
		if !inString && !inComment && i+1 < len(content) {
			if content[i:i+2] == "--" {
				inComment = true
				commentType = "--"
				i++ // Пропускаем следующий символ
				continue
			} else if content[i:i+2] == "/*" {
				inComment = true
				commentType = "/*"
				i++ // Пропускаем следующий символ
				continue
			}
		}

		// Выход из комментариев
		if inComment {
			if commentType == "--" && char == '\n' {
				inComment = false
			} else if commentType == "/*" && i+1 < len(content) && content[i:i+2] == "*/" {
				inComment = false
				i++ // Пропускаем следующий символ
			}
			continue
		}

		// Обработка строк
		if char == '\'' || char == '"' {
			if !inString {
				inString = true
				stringChar = char
			} else if stringChar == char {
				// Проверяем, не экранирована ли кавычка
				if i > 0 && content[i-1] != '\\' {
					inString = false
				}
			}
		}

		// Разделение запросов
		if char == ';' && !inString {
			query := strings.TrimSpace(currentQuery.String())
			if query != "" {
				queries = append(queries, query)
			}
			currentQuery.Reset()
		} else {
			currentQuery.WriteByte(char)
		}
	}

	// Добавляем последний запрос, если он есть
	finalQuery := strings.TrimSpace(currentQuery.String())
	if finalQuery != "" {
		queries = append(queries, finalQuery)
	}

	return queries
}

// extractMigrationNumber извлекает номер миграции из имени файла
func extractMigrationNumber(filename string) int {
	re := regexp.MustCompile(`^(\d+)_`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		var num int
		fmt.Sscanf(matches[1], "%d", &num)
		return num
	}
	return -1
}

// isMigrationApplied проверяет, была ли применена миграция
func isMigrationApplied(db *sql.DB, version string) bool {
	// Создаем таблицу для отслеживания миграций, если её нет
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Warning: failed to create schema_migrations table: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		// Если таблицы нет, то считаем что миграция не применена
		if strings.Contains(err.Error(), "does not exist") {
			return false
		}
		log.Printf("Warning: failed to check migration status: %v", err)
		return false
	}

	return count > 0
}
