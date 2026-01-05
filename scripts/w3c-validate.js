/**
 * W3C HTML Validator Script
 * 
 * This script validates HTML files using the W3C Validator API.
 * It reads HTML files from the frontend/templates directory and sends them to the W3C API.
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import fetch from 'node-fetch';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const templatesDir = path.join(__dirname, '..', 'frontend', 'templates');

// W3C Validator API endpoint
const W3C_VALIDATOR_URL = 'https://validator.w3.org/nu/?out=json';

/**
 * Find all HTML files in a directory recursively
 * @param {string} dir - Directory to search
 * @param {Array} fileList - Accumulator for found files
 * @returns {Array} - List of HTML file paths
 */
function findHtmlFiles(dir, fileList = []) {
    const files = fs.readdirSync(dir);

    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);

        if (stat.isDirectory()) {
            findHtmlFiles(filePath, fileList);
        } else if (path.extname(file) === '.html') {
            fileList.push(filePath);
        }
    });

    return fileList;
}

/**
 * Validate an HTML file using the W3C Validator API
 * @param {string} filePath - Path to the HTML file
 */
async function validateHtmlFile(filePath) {
    try {
        const html = fs.readFileSync(filePath, 'utf8');

        // Add DOCTYPE if missing (for template fragments)
        let htmlContent = html;
        if (!htmlContent.trim().startsWith('<!DOCTYPE')) {
            htmlContent = `<!DOCTYPE html>\n<html>\n<head>\n<title>Test</title>\n</head>\n<body>\n${htmlContent}\n</body>\n</html>`;
        }

        const response = await fetch(W3C_VALIDATOR_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'text/html; charset=utf-8',
                'User-Agent': 'Mozilla/5.0 (Node.js W3C Validator Script)'
            },
            body: htmlContent
        });

        const result = await response.json();

        // Process validation results
        const errors = result.messages.filter(msg => msg.type === 'error');
        const warnings = result.messages.filter(msg => msg.type === 'info' || msg.type === 'warning');

        const relativePath = path.relative(path.join(__dirname, '..'), filePath);

        if (errors.length > 0 || warnings.length > 0) {
            console.log(`\n\x1b[1m${relativePath}\x1b[0m`);

            if (errors.length > 0) {
                console.log(`\x1b[31mErrors (${errors.length}):\x1b[0m`);
                errors.forEach(error => {
                    console.log(`  Line ${error.lastLine}: ${error.message}`);
                });
            }

            if (warnings.length > 0) {
                console.log(`\x1b[33mWarnings (${warnings.length}):\x1b[0m`);
                warnings.forEach(warning => {
                    console.log(`  Line ${warning.lastLine}: ${warning.message}`);
                });
            }
        } else {
            console.log(`\x1b[32m✓ ${relativePath} - No issues found\x1b[0m`);
        }

    } catch (error) {
        console.error(`\x1b[31mError validating ${filePath}: ${error.message}\x1b[0m`);
    }
}

/**
 * Main function to validate all HTML files
 */
async function validateAllHtmlFiles() {
    console.log('\x1b[1mStarting W3C HTML validation...\x1b[0m');

    try {
        const htmlFiles = findHtmlFiles(templatesDir);
        console.log(`Found ${htmlFiles.length} HTML files to validate.\n`);

        for (const file of htmlFiles) {
            await validateHtmlFile(file);
        }

        console.log('\n\x1b[1mW3C HTML validation completed.\x1b[0m');
    } catch (error) {
        console.error(`\x1b[31mError during validation: ${error.message}\x1b[0m`);
        process.exit(1);
    }
}

// Run the validation
validateAllHtmlFiles();