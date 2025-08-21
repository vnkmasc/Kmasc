// components/tiny-text-edit.tsx
'use client'

import { forwardRef, useImperativeHandle, useRef } from 'react'
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
    getContent: () => editorRef.current?.getContent?.() ?? ''
  }))

  return (
    <Editor
      tinymceScriptSrc='/assets/libs/tinymce/tinymce.min.js'
      licenseKey='gpl'
      onInit={(_, editor) => (editorRef.current = editor)}
      value={value}
      onEditorChange={(content) => onChange?.(content)}
      init={{
        base_url: '/assets/libs/tinymce',
        height: 700,
        menubar: false,
        plugins: [
          'advlist',
          'autolink',
          'lists',
          'link',
          'image',
          'charmap',
          'preview',
          'anchor',
          'searchreplace',
          'visualblocks',
          'code',
          'fullscreen',
          'insertdatetime',
          'media',
          'table',
          'help',
          'wordcount'
        ],
        toolbar:
          'undo redo | blocks | bold italic forecolor | alignleft aligncenter alignright alignjustify | ' +
          'bullist numlist outdent indent | removeformat | code preview fullscreen | help',
        content_style: 'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }'
      }}
    />
  )
})

TinyTextEdit.displayName = 'TinyTextEdit'
export default TinyTextEdit
