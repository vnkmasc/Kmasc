// next.config.mjs
import path from 'path'
import CopyPlugin from 'copy-webpack-plugin'
import { NextConfig } from 'next'

/** @type {import('next').NextConfig} */
const nextConfig: NextConfig = {
  webpack: (config) => {
    config.plugins.push(
      new CopyPlugin({
        patterns: [
          {
            from: path.join(process.cwd(), 'node_modules/tinymce/skins'),
            to: path.join(process.cwd(), 'public/assets/libs/tinymce/skins')
          }
        ]
      })
    )
    return config
  }
}

export default nextConfig
