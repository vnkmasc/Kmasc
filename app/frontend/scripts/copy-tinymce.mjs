import fse from 'fs-extra'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const projectRoot = path.join(__dirname, '..')

const src = path.join(projectRoot, 'node_modules', 'tinymce')
const dest = path.join(projectRoot, 'public', 'assets', 'libs', 'tinymce')

await fse.emptyDir(dest)
await fse.copy(src, dest, { overwrite: true })
