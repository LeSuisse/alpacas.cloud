import {defineConfig} from 'vite'
import copy from "rollup-plugin-copy";

export default defineConfig({
    build: {
        commonjsOptions: {
            include: []
        }
    },
    optimizeDeps: {
        disabled: false,
    },
    plugins: [
        copy({
            targets: [{src: "openapi.json", dest: "dist/"}],
            hook: 'writeBundle'
        })
    ]
});
