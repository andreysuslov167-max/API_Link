package handlers

import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "fmt"

    "github.com/gofiber/fiber/v3"

    "project-go/db"
    "project-go/models"
)

type LinkHandler struct{}

func NewLinkHandler() *LinkHandler {
    return &LinkHandler{}
}

func (h *LinkHandler) CreateLink(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int)

    var req models.CreateLinkRequest
    if err := c.Bind().JSON(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if req.URL == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "URL is required",
        })
    }

    shortCode := generateShortCode()

    var link models.Link
    err := database.DB.QueryRow(
        "INSERT INTO links (user_id, original_url, short_code) VALUES (?, ?, ?) RETURNING id, user_id, original_url, short_code, clicks, created_at",
        userID, req.URL, shortCode,
    ).Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Clicks, &link.CreatedAt)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error creating link",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "link":      link,
        "short_url": fmt.Sprintf("%s/%s", c.Hostname(), link.ShortCode),
    })
}

func (h *LinkHandler) GetAllLinks(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int)

    rows, err := database.DB.Query(
        "SELECT id, user_id, original_url, short_code, clicks, created_at FROM links WHERE user_id = ? ORDER BY created_at DESC",
        userID,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
        })
    }
    defer rows.Close()

    var links []models.Link
    for rows.Next() {
        var link models.Link
        if err := rows.Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Clicks, &link.CreatedAt); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Error scanning links",
            })
        }
        links = append(links, link)
    }

    if links == nil {
        links = []models.Link{}
    }

    return c.JSON(fiber.Map{
        "links": links,
        "total": len(links),
    })
}

func (h *LinkHandler) GetLink(c fiber.Ctx) error {
    id := c.Params("id")
    userID := c.Locals("user_id").(int)

    var link models.Link
    err := database.DB.QueryRow(
        "SELECT id, user_id, original_url, short_code, clicks, created_at FROM links WHERE id = ? AND user_id = ?",
        id, userID,
    ).Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Clicks, &link.CreatedAt)

    if err == sql.ErrNoRows {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Link not found",
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
        })
    }

    return c.JSON(link)
}

func (h *LinkHandler) DeleteLink(c fiber.Ctx) error {
    id := c.Params("id")
    userID := c.Locals("user_id").(int)

    result, err := database.DB.Exec(
        "DELETE FROM links WHERE id = ? AND user_id = ?",
        id, userID,
    )

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
        })
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Link not found",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Link deleted successfully",
    })
}

func (h *LinkHandler) Redirect(c fiber.Ctx) error {
    code := c.Params("code")

    var link models.Link
    err := database.DB.QueryRow(
        "SELECT id, original_url FROM links WHERE short_code = ?",
        code,
    ).Scan(&link.ID, &link.OriginalURL)

    if err == sql.ErrNoRows {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Link not found",
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
        })
    }

    go func() {
        database.DB.Exec(
            "UPDATE links SET clicks = clicks + 1 WHERE id = ?",
            link.ID,
        )
    }()

    return c.Redirect().To(link.OriginalURL)
}

func (h *LinkHandler) GetStats(c fiber.Ctx) error {
    id := c.Params("id")
    userID := c.Locals("user_id").(int)

    var link models.Link
    err := database.DB.QueryRow(
        "SELECT id, user_id, original_url, short_code, clicks, created_at FROM links WHERE id = ? AND user_id = ?",
        id, userID,
    ).Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Clicks, &link.CreatedAt)

    if err == sql.ErrNoRows {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Link not found",
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
        })
    }

    stats := models.LinkStats{
        Link:      link,
        Redirects: link.Clicks,
        ShortURL:  fmt.Sprintf("%s/%s", c.Hostname(), link.ShortCode),
    }

    return c.JSON(stats)
}

func generateShortCode() string {
    b := make([]byte, 6)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)[:8]
}