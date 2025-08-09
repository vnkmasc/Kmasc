import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import { ThemeProvider } from '@/components/providers/theme-provider'
import SWRConfig from '@/components/providers/swr-config'
import './globals.css'
import { Toaster } from 'sonner'
import { TooltipProvider } from '@/components/ui/tooltip'

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin']
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin']
})

export const metadata: Metadata = {
  title: 'Kmasc',
  description: 'Giải pháp quản lý văn bằng chứng chỉ ứng dụng Blockchain.',
  icons: {
    icon: '/assets/images/logoKMA.png'
  }
}

export default async function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang={'vi'} suppressHydrationWarning>
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <ThemeProvider
          attribute='class'
          defaultTheme='system'
          enableSystem
          disableTransitionOnChange
          storageKey='theme'
        >
          <TooltipProvider>
            <SWRConfig>
              {children}
              <Toaster expand={true} />
            </SWRConfig>
          </TooltipProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}
