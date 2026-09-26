import { cp, mkdir } from 'node:fs/promises';

const root = new URL('./', import.meta.url);
await mkdir(new URL('dist/assets/', root), { recursive: true });
await cp(new URL('index.html', root), new URL('dist/index.html', root));
await cp(new URL('assets/', root), new URL('dist/assets/', root), { recursive: true });
