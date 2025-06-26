'use client'

import { useTheme } from 'next-themes'
import { Button } from '../ui/button'
import { MoonIcon, SunIcon } from 'lucide-react'
import { useEffect, useState } from 'react'

const ThemeSwitch = () => {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return (
      <Button variant='ghost' size='icon'>
        <div className='h-4 w-4' />
      </Button>
    )
  }

  return (
    <Button variant='ghost' size='icon' onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>
      {theme === 'light' ? <SunIcon className='h-4 w-4' /> : <MoonIcon className='h-4 w-4' />}
    </Button>
  )
}

export default ThemeSwitch
