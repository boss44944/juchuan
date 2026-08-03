import { cp, mkdir, rename, rm } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const frontendRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const source = path.join(frontendRoot, 'dist')
const target = path.resolve(frontendRoot, '..', 'backend', 'static')
const stagedTarget = `${target}.sync-${process.pid}`

await rm(stagedTarget, { recursive: true, force: true })
await mkdir(stagedTarget, { recursive: true })
await cp(source, stagedTarget, { recursive: true })
await rm(target, { recursive: true, force: true })
await rename(stagedTarget, target)

console.log(`Synced frontend build to ${target}`)
