# MSDS 431 Assignment 10 — GenAI Conversation Log

## Tool Used

**OpenAI ChatGPT**

## Purpose

ChatGPT was used throughout Assignment 10 as a development and troubleshooting assistant while building a movie database application using Go, Python, Polars, and DuckDB.

This document summarizes the main areas of AI assistance. It is a summary of the development conversations, not a verbatim transcript.

## Areas of Assistance

### 1. Project Planning and Database Design

ChatGPT assisted with organizing the assignment structure and planning the integration of Python, Polars, DuckDB, and Go. Assistance included reviewing the database schema, identifying relationships between movies, genres, actors, and roles, and planning the ingestion and validation workflow.

### 2. Python Data Ingestion and Validation

ChatGPT provided guidance on using Polars to load the source CSV files and DuckDB to create and populate the database. It also assisted with validation queries and interpreting record counts, missing role values, duplicate identifiers, and orphaned records.

### 3. Go Application Development

ChatGPT assisted with developing the Go command-line application and its integration with DuckDB. This included database initialization, SQL query implementation, application organization, and development of the interactive menu and its features.

### 4. Movie Database Features

ChatGPT provided implementation guidance and troubleshooting support for movie searching, genre exploration, actor searches, top-ranked movie listings, genre analysis, movie details, and personal collection management.

### 5. Testing and Software Quality

ChatGPT assisted with developing database and collection workflow tests, diagnosing test failures, improving test coverage, and reviewing validation results. This included testing database constraints, duplicate collection entries, invalid ratings, and collection operations.

### 6. Benchmarking and Performance Analysis

ChatGPT assisted with creating and running Go benchmarks for movie search, genre analysis, and actor search. It also helped interpret benchmark execution times, memory usage, and allocation counts.

### 7. Environment Configuration and Troubleshooting

ChatGPT helped troubleshoot Windows-specific development issues, including Go CGO configuration, GCC compiler selection, DuckDB driver integration, and application build errors.

### 8. Documentation

ChatGPT assisted with organizing the project documentation, writing setup and execution instructions, documenting test coverage and benchmark results, and preparing the README and GenAI disclosure.

## Student Contributions

The student directed the implementation process and was responsible for executing the development workflow. This included configuring the local development environment, running Python ingestion and validation scripts, integrating the Go application, testing its features, reviewing database results, and executing Go tests and benchmarks.

The student reviewed and evaluated AI-generated suggestions, made implementation decisions, and verified the final application behavior and reported results.

## Verification

All reported database counts, validation results, test outcomes, coverage measurements, and benchmark measurements were obtained by executing the corresponding scripts and commands in the local development environment. AI-generated explanations and suggestions were reviewed rather than accepted without verification.