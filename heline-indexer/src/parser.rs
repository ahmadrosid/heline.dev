use std::fs;
use std::path::PathBuf;

// TODO: Use this enum to map the file extension
#[allow(dead_code)]
enum Language {
    Raw,
    Shell,
    Javascript,
    Go,
    Groovy,
    Typescript,
    C,
    Cpp,
    HTML,
    Haskell,
    Java,
    Kotlin,
    Lua,
    Markdown,
    Nim,
    Python,
    PHP,
    Ruby,
    Rust,
    Toml,
    CSharp,
    Yaml,
    Dart,
    Diff,
    Json,
    ProtocolBuffers,
    Gemfile,
    Dockerfile,
}

pub fn read_file(file_path: &PathBuf) -> Result<(Vec<char>, &str), String> {
    let path = file_path.to_str().unwrap();
    if let Some(name) = file_path.file_name().unwrap().to_str() {
        match name {
            "package-lock.json" => return Err(format!("Ignore file: '{}'!", path)),
            "yarn.lock" => return Err(format!("Ignore file: '{}'!", path)),
            _ => {}
        }
    }

    if let Ok(source) = fs::read_to_string(path) {
        if source.len() == 0 {
            return Err(format!("Failed to read file: '{}'!", path));
        }

        if source.len() <= 12 {
            return Err(format!("Source code to short: {}", source));
        }

        let input = source.chars().collect();
        let lang = match file_path.extension() {
            Some(ext) => match ext.to_str().unwrap_or("Raw") {
                "sh" | "zsh" | "bash" => "Shell",
                "js" => "JavaScript",
                "go" => "Go",
                "groovy" => "Groovy",
                "ts" | "tsx" => "TypeScript",
                "c" | "h" => "C",
                "cpp" | "c++" => "C++",
                "html" => "HTML",
                "hs" => "Haskell",
                "java" => "Java",
                "kt" => "Kotlin",
                "lua" => "Lua",
                "md" | "adoc" => "Markdown",
                "nim" => "Nim",
                "py" => "Python",
                "php" => "PHP",
                "ru" | "rb" | "podspec" => "Ruby",
                "rs" => "Rust",
                "toml" => "TOML",
                "cs" => "C#",
                "yml" | "yaml" => "YAML",
                "dart" => "Dart",
                "patch" => "Diff",
                "json" => "JSON",
                "proto" => "Protocol Buffer",
                "lock" => match file_path.file_name().unwrap().to_str().unwrap() {
                    "Cargo.lock" => "TOML",
                    "Gemfile.lock" => "Gemfile",
                    "yarn.lock" => "YAML",
                    _ => "Raw",
                },
                _ => ext.to_str().unwrap(),
            },
            _ => parse_file_name(file_path.file_name().unwrap().to_str().unwrap()),
        };
        Ok((input, lang))
    } else {
        Err(format!("Failed to read file: '{}'!", path))
    }
}

fn parse_file_name(file: &str) -> &str {
    return match file {
        "Jenkinsfile" => "Groovy",
        "Dockerfile" => "Dockerfile",
        "Makefile" => "Makefile",
        "Gemfile" => "Gemfile",
        "Rakefile" => "Rakefile",
        _ => "Raw",
    };
}

pub fn render_html(input: Vec<char>, lang: &str) -> String {
    return match lang {
        "Shell" => hl_core::render_html(input, "bash"),
        "C" => hl_core::render_html(input, "c"),
        "C++" => hl_core::render_html(input, "cpp"),
        "Clojure" => hl_core::render_html(input, "clojure"),
        "CSS" => hl_core::render_html(input, "css"),
        "CUDA" => hl_core::render_html(input, "cuda"),
        "Dart" => hl_core::render_html(input, "dart"),
        "edn" => hl_core::render_html(input, "edn"),
        "Go" => hl_core::render_html(input, "go"),
        "Groovy" => hl_core::render_html(input, "groovy"),
        "Haskell" => hl_core::render_html(input, "haskell"),
        "HTML" => hl_core::render_html(input, "html"),
        "Ruby" | "Rakefile" | "Gemfile" => hl_core::render_html(input, "ruby"),
        "Rust" => hl_core::render_html(input, "rust"),
        "C#" => hl_core::render_html(input, "cs"),
        "Java" => hl_core::render_html(input, "java"),
        "JavaScript" => hl_core::render_html(input, "javascript"),
        "JSON" => hl_core::render_html(input, "json"),
        "Kotlin" => hl_core::render_html(input, "kotlin"),
        "Lua" => hl_core::render_html(input, "lua"),
        "Makefile" => hl_core::render_html(input, "makefile"),
        "Markdown" => hl_core::render_html(input, "markdown"),
        "Nim" => hl_core::render_html(input, "nim"),
        "PHP" => hl_core::render_html(input, "php"),
        "Python" => hl_core::render_html(input, "python"),
        "TOML" => hl_core::render_html(input, "toml"),
        "TypeScript" => hl_core::render_html(input, "typescript"),
        "YAML" => hl_core::render_html(input, "yaml"),
        "Protocol Buffer" => hl_core::render_html(input, "proto"),
        _ => {
            let mark_bash = String::from("#!/bin/sh");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "bash");
                }
            }

            let mark_bash = String::from("#!/bin/bash");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "bash");
                }
            }

            let mark_bash = String::from("#!/usr/bin/env bash");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "bash");
                }
            }

            let mark_bash = String::from("#!/usr/bin/env sh");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "bash");
                }
            }

            let mark_bash = String::from("#!/usr/bin/env zsh");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "bash");
                }
            }

            let mark_bash = String::from("#!/usr/bin/env ruby");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "ruby");
                }
            }

            let mark_bash = String::from("#!/usr/bin/env php");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "php");
                }
            }

            let mark_bash = String::from("@ruby");
            if input.len() > mark_bash.len() {
                let result: String = input[0..mark_bash.len()].iter().collect();
                if result == mark_bash {
                    return hl_core::render_html(input, "ruby");
                }
            }

            return hl_core::render_html(input, "raw");
        }
    };
}
