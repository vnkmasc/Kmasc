'use client'

import { useRef, useImperativeHandle, forwardRef } from 'react'
import { Editor } from '@tinymce/tinymce-react'

interface TinyTextEditProps {
  value?: string
  // eslint-disable-next-line no-unused-vars
  onChange?: (content: string) => void
}

export interface TinyTextEditRef {
  getContent: () => string
}

const TinyTextEdit = forwardRef<TinyTextEditRef, TinyTextEditProps>(({ value, onChange }, ref) => {
  const editorRef = useRef<any>(null)

  useImperativeHandle(ref, () => ({
    getContent: () => {
      if (editorRef.current) {
        return editorRef.current.getContent()
      }
      return ''
    }
  }))

  return (
    <Editor
      apiKey={process.env.NEXT_PUBLIC_TINYMCE_API_KEY}
      onInit={(_evt, editor) => (editorRef.current = editor)}
      value={value}
      onEditorChange={(content) => {
        if (onChange) {
          onChange(content)
        }
      }}
      init={{
        height: 700,
        menubar: false,
        plugins: [
          'advlist',
          'autolink',
          'lists',
          'charmap',
          'preview',
          'anchor',
          'searchreplace',
          'visualblocks',
          'fullscreen',
          'insertdatetime',
          'media',
          'code',
          'help',
          'wordcount'
        ],
        toolbar:
          'undo redo | blocks | ' +
          'bold italic forecolor | alignleft aligncenter ' +
          'alignright alignjustify | bullist numlist outdent indent | ' +
          'removeformat | help',
        content_style: 'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }',
        skin_url: '/assets/libs/tinymce/skins/ui/oxide',
        content_css: '/assets/libs/tinymce/skins/content/default/content.min.css'
      }}
    />
  )
})

TinyTextEdit.displayName = 'TinyTextEdit'

export default TinyTextEdit
