import Footer from '@/components/common/footer'
import Header from '@/components/common/header'
import { DataProvider } from '@/components/providers/data-provider'

interface Props {
  children: React.ReactNode
}

const EducationAdminLayout: React.FC<Props> = ({ children }) => {
  return (
    <DataProvider>
      <main className='flex h-screen flex-col'>
        <Header role='university_admin' />
        <div className='container mt-16 flex-1 py-6'>{children}</div>
        <Footer />
      </main>
    </DataProvider>
  )
}

export default EducationAdminLayout
