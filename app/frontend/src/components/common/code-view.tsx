'use client'

import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'
import { codeToHtml } from 'shiki'

interface Props {
  code: string
}

const CodeView: React.FC<Props> = ({ code }) => {
  const [html, setHtml] = useState<string>('')
  const { theme } = useTheme()

  useEffect(() => {
    const fetchHtml = async () => {
      const html = await codeToHtml(code, {
        lang: 'html',
        themes: {
          light: theme === 'dark' ? 'dracula' : 'github-light',
          dark: 'nord'
        }
      })
      setHtml(html)
    }
    fetchHtml()
  }, [code, theme])

  return (
    <div
      className='overflox-x-scroll overflow-y-scroll rounded-lg border border-gray-200 p-4 dark:border-gray-700'
      dangerouslySetInnerHTML={{ __html: html }}
    />
  )
}

export default CodeView
