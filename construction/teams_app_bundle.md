# Teams App Bundling and Deployment

This document describes how to package the Teams App for local development testing and eventual production deployment.

## Overview

Packaging the Teams App involves creating a ZIP archive containing the application's metadata (manifest) and visual identity (icons). This package is then uploaded to Microsoft Teams.

## Prerequisites

1.  **Public URL**: Since Microsoft Teams needs to load the application over HTTPS, you must expose your local development server (`http://localhost:5173`) using a tunnel like `ngrok` or `dev-tunnels`.
2.  **App Icons**: Ensure `color.png` (192x192) and `outline.png` (32x32) exist in `src/frontend/teams-app/public/`.

## Packaging Process

We use a helper script `package-app.sh` located in the root directory to automate the bundling.

### Using the Packaging Script

Run the script from the project root with your public HTTPS URL:

```bash
sh package-app.sh <https-tunnel-url>
```

Example:
```bash
sh package-app.sh https://random-id.ngrok-free.app
```

### What the script does:
1.  **Manifest Update**: Creates a temporary `manifest.json` replacing placeholders (`https://YOUR_NGROK_URL` and `YOUR_NGROK_DOMAIN`) with your actual tunnel URL and domain.
2.  **Validation**: Verifies that the required icons exist.
3.  **Bundling**: Creates an `app.zip` file containing:
    *   `manifest.json`
    *   `color.png`
    *   `outline.png`
4.  **Output**: Saves `app.zip` in the project root.

## Uploading to Microsoft Teams

1.  Open **Microsoft Teams** (desktop or web).
2.  Navigate to the **Apps** section on the left sidebar.
3.  Select **Manage your apps** at the bottom.
4.  Click **Upload an app** -> **Upload a custom app**.
5.  Select the generated `app.zip` file.
6.  Click **Add** to install the app for yourself.

## Troubleshooting

-   **Manifest Errors**: If Teams rejects the ZIP, ensure the `manifest.json` version and ID match the requirements in `src/frontend/teams-app/public/manifest.json`.
-   **HTTPS Only**: Microsoft Teams will not load `http` content. Ensure your tunnel provides a valid `https` URL.
-   **Valid Domains**: Any external API or resource the app calls must be listed in the `validDomains` section of `manifest.json`. The packaging script automatically adds the tunnel domain.
