import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import fs from 'fs';
import path from 'path';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), ''); 

    return {
        plugins: [react()],
        server: {
            https: {
                key: fs.readFileSync(path.resolve(__dirname, env.TLS_KEY_PATH)),
                cert: fs.readFileSync(path.resolve(__dirname, env.TLS_CERT_PATH)),
            },
            host: env.SERVER_NAME,
            port: 5173,
        }
    }
})
