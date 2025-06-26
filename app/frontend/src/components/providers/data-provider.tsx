'use client'

import { getFacultyList } from '@/lib/api/faculty'
import React, { createContext, useContext, ReactNode } from 'react'
import useSWR, { mutate } from 'swr'

type Faculty = {
  code: string
  name: string
  id: string
}

// Create the context with a default value
const DataContext = createContext<{
  facultyList: Faculty[]
}>({
  facultyList: []
})

// Provider component interface
interface DataProviderProps {
  children: ReactNode
}

// Provider component
export const DataProvider: React.FC<DataProviderProps> = ({ children }) => {
  const queryFaculty = useSWR('faculty-list', getFacultyList, {
    refreshInterval: 1000 * 60 * 5
  })

  return <DataContext.Provider value={{ facultyList: queryFaculty.data || [] }}>{children}</DataContext.Provider>
}

// Custom hook to use the context
export const UseData = (): {
  facultyList: Faculty[]
} => {
  const context = useContext(DataContext)

  if (context === undefined) {
    throw new Error('useData must be used within a DataProvider')
  }

  return context
}

export const UseRefetchFacultyList = () => {
  return mutate('faculty-list')
}

// Export the context for direct use if needed
export { DataContext }
