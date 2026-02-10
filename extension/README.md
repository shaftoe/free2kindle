# savetoink Firefox Extension

A simple Firefox extension to send web articles to your Kindle device.

## Features

- **Click to Send**: Click the extension icon to send the current page to your Kindle
- **Context Menu**: Right-click on any link to send it to your Kindle
- **Secure Settings**: Store your API key and URL securely using Chrome storage API

## Installation

1. Open Firefox and navigate to `about:debugging#/runtime/this-firefox`
2. Click "Load Temporary Add-on"
3. Select the `manifest.json` file in this directory
4. Configure your API settings by clicking the extension icon and selecting "Configure Settings"

## Configuration

Before using the extension, you need to configure:

1. **API URL**: The URL of your deployed savetoink Lambda function
2. **API Key**: Your secret API key for authentication

To configure settings:
- Click the extension icon
- Click "Configure Settings"
- Enter your API URL and API key
- Click "Save Settings"

The extension will verify your settings by making a health check request to your API.

## Usage

### Send Current Page
1. Navigate to any article or web page
2. Click the savetoink extension icon
3. Click "Send to Kindle"
4. Wait for confirmation message

### Send a Link
1. Right-click on any link on a web page
2. Select "Send to Kindle" from the context menu
3. Wait for confirmation message

## API Compatibility

This extension works with the savetoink backend API:

- `POST /v1/articles` - Queue article for Kindle delivery
 - `GET /v1/health` - Health check endpoint

Authentication is done via the `Authorization` header with `Bearer` prefix.

## Security

- API keys are stored using `chrome.storage.local`, which is not encrypted but isolated per extension
- For production use, consider using `chrome.storage.session` for ephemeral storage
- The extension only communicates with the configured API URL

## Development

The extension uses Manifest V3 and is compatible with Firefox and Chrome.

File structure:
- `manifest.json` - Extension manifest
- `popup.html` - Popup UI for icon click
- `popup.js` - Popup logic for sending current tab
- `background.js` - Background script for context menu
- `options.html` - Settings page UI
- `options.js` - Settings page logic
- `icons/` - Extension icons (you need to add your own icons)

## Icons

You need to add the following icon files to the `icons/` directory:
- `icon-16.png` (16x16px)
- `icon-32.png` (32x32px)
- `icon-48.png` (48x48px)
- `icon-128.png` (128x128px)

You can use any PNG images for these icons. A simple book or Kindle icon would work well.

## License

Same as the main savetoink project.
