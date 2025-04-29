package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Calculator Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			return mcp.NewToolResultError("未対応の機能です"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	s.AddTool(
		mcp.NewTool(
			"uuid",
			mcp.WithDescription("generate uuid"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uuid, err := uuid.NewUUID()
			if err != nil {
				return mcp.NewToolResultError("uuidの生成に失敗しました"), err
			}
			return mcp.NewToolResultText(uuid.String()), nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			"docs://readme",
			"Project README",
			mcp.WithResourceDescription("The project's README file"),
			mcp.WithMIMEType("text/markdown"),
		),
		func(
			ctx context.Context,
			request mcp.ReadResourceRequest,
		) ([]mcp.ResourceContents, error) {
			content := `
				# MCP demo
				- calculator
				- generate uuid
				`

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "docs://content",
					MIMEType: "text/markdown",
					Text:     string(content),
				},
			}, nil
		})

	s.AddPrompt(
		mcp.NewPrompt("query_builder",
			mcp.WithPromptDescription("SQL query builder assistance"),
			mcp.WithArgument("table",
				mcp.ArgumentDescription("Name of the table to query"),
				mcp.RequiredArgument(),
			),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			tableName := request.Params.Arguments["table"]
			if tableName == "" {
				return nil, fmt.Errorf("table name is required")
			}

			return mcp.NewGetPromptResult(
				"SQL query builder assistance",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(
						mcp.RoleAssistant,
						mcp.NewTextContent("You are a SQL expert. Help construct efficient and safe queries."),
					),
					mcp.NewPromptMessage(
						mcp.RoleAssistant,
						mcp.NewEmbeddedResource(mcp.TextResourceContents{
							URI:      fmt.Sprintf("db://schema/%s", tableName),
							MIMEType: "application/json",
						}),
					),
				},
			), nil
		})

	// サーバー起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
