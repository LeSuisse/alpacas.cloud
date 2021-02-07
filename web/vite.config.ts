import {defineConfig} from 'vite'
import copy from "rollup-plugin-copy";

export default defineConfig({
    plugins: [
        copy({
            targets: [{src: "openapi.json", dest: "dist/"}],
            hook: 'writeBundle'
        })
    ]
});
