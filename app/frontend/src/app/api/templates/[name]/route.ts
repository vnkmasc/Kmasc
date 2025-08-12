import { NextResponse } from 'next/server'
import fs from 'node:fs/promises'
import path from 'node:path'

export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

const ALLOWED_TEMPLATES = new Set(['v1-degree', 'v2-degree', 'v3-degree'])

export async function GET(_req: Request, context: { params: Promise<{ name: string }> }) {
  const { name } = await context.params

  if (!ALLOWED_TEMPLATES.has(name)) {
    return new NextResponse('Not found', { status: 404 })
  }

  const fileName = `${name}.html`

  // Try multiple paths for different deployment scenarios
  const possiblePaths = [
    path.join(process.cwd(), 'public', 'templates', fileName),
    path.join(process.cwd(), 'src', 'templates', fileName),
    path.join(process.cwd(), 'app', 'frontend', 'public', 'templates', fileName),
    path.join(process.cwd(), 'app', 'frontend', 'src', 'templates', fileName)
  ]

  for (const filePath of possiblePaths) {
    try {
      const content = await fs.readFile(filePath, 'utf8')
      return new NextResponse(content, {
        status: 200,
        headers: {
          'content-type': 'text/html; charset=utf-8'
        }
      })
    } catch {
      continue
    }
  }
  return new NextResponse('Template not found', { status: 404 })
}
