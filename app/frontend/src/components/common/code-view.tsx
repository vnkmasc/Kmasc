'use client'

import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'
import { codeToHtml } from 'shiki'
import { Copy, Check } from 'lucide-react'

interface Props {
  code: string
}

const CodeView: React.FC<Props> = ({ code }) => {
  const [html, setHtml] = useState<string>('')
  const [copied, setCopied] = useState(false)
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

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => {
        setCopied(false)
      }, 200)
    } catch (err) {
      console.error('Failed to copy text: ', err)
    }
  }

  return (
    <div className='relative overflow-auto rounded-lg border border-gray-200 p-4 dark:border-gray-700'>
      <button
        onClick={handleCopy}
        className='absolute right-2 top-2 rounded p-1 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800'
        title={copied ? 'Copied!' : 'Copy code'}
      >
        {copied ? <Check className='text-green-500' /> : <Copy />}
      </button>
      <div dangerouslySetInnerHTML={{ __html: html }} />
    </div>
  )
}

export default CodeView
