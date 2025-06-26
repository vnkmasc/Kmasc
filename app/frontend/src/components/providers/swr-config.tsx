'use client'

import apiService from '@/lib/api/root'
import { SWRConfig as SWRConfigNext } from 'swr'

const SWRConfig = ({ children }: { children: React.ReactNode }) => {
  return (
    <SWRConfigNext
      value={{
        loadingTimeout: 5000,
        shouldRetryOnError: false,
        revalidateOnFocus: false,
        fetcher: (url) => apiService('GET', url)
      }}
    >
      {children}
    </SWRConfigNext>
  )
}

export default SWRConfig
