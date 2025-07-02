import Footer from '@/components/common/footer'
import Header from '@/components/common/header'

interface Props {
  children: React.ReactNode
}

const StudentLayout: React.FC<Props> = ({ children }) => {
  return (
    <main className='flex h-screen flex-col'>
      <Header role='student' />
      <div className='container mt-16 flex-1 py-6'>{children}</div>
      <Footer />
    </main>
  )
}

export default StudentLayout
