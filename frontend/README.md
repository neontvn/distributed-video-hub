This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Environment Variables

This application requires the following environment variable:

- `NEXT_PUBLIC_API_URL`: The URL of your backend API server (e.g., your GCP server endpoint)

For local development, create a `.env.local` file:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Deploy on Vercel

### Prerequisites
1. Make sure your Go backend server is running on GCP and accessible via HTTPS
2. Note down your GCP server's public URL

### Deployment Steps

1. **Install Vercel CLI** (if not already installed):
   ```bash
   npm i -g vercel
   ```

2. **Login to Vercel**:
   ```bash
   vercel login
   ```

3. **Navigate to the Next.js directory**:
   ```bash
   cd triton-tube-nextjs
   ```

4. **Deploy to Vercel**:
   ```bash
   vercel
   ```

5. **Set Environment Variables**:
   - Go to your Vercel dashboard
   - Select your project
   - Go to Settings → Environment Variables
   - Add the following variable:
     - **Name**: `NEXT_PUBLIC_API_URL`
     - **Value**: Your GCP server URL (e.g., `https://your-gcp-server.com`)
     - **Environment**: Production, Preview, Development

6. **Redeploy** (if needed):
   ```bash
   vercel --prod
   ```

### Alternative: Deploy via GitHub

1. Push your code to a GitHub repository
2. Connect your repository to Vercel
3. Set the environment variables in the Vercel dashboard
4. Vercel will automatically deploy on every push

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/vercel.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
