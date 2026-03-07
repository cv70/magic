import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import { Button, Space, Tabs, Tooltip, message } from 'antd'
import {
  BoldOutlined,
  ItalicOutlined,
  StrikethroughOutlined,
  UnorderedListOutlined,
  OrderedListOutlined,
  ClearOutlined,
  CopyOutlined,
  UndoOutlined,
  RedoOutlined,
} from '@ant-design/icons'
import './TipTapEditor.css'

interface TipTapEditorProps {
  value: string
  onChange: (content: string) => void
  placeholder?: string
  height?: string
}

export function TipTapEditor({ value, onChange, placeholder = '输入内容...', height = '400px' }: TipTapEditorProps) {
  const editor = useEditor({
    extensions: [StarterKit],
    content: value || '',
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML())
    },
  })

  if (!editor) {
    return <div>加载编辑器中...</div>
  }

  const handleFormatClick = (command: () => void) => {
    command()
    editor.view.focus()
  }

  const copyMarkdown = () => {
    const text = editor.getJSON()
    navigator.clipboard.writeText(JSON.stringify(text, null, 2))
    message.success('已复制 JSON 内容')
  }

  const clearContent = () => {
    if (window.confirm('确定清空内容？')) {
      editor.commands.clearContent()
      onChange('')
    }
  }

  return (
    <div className="tiptap-editor-wrapper">
      <Tabs
        defaultActiveKey="wysiwyg"
        items={[
          {
            key: 'wysiwyg',
            label: '富文本编辑',
            children: (
              <div>
                {/* 工具栏 */}
                <div className="tiptap-toolbar">
                  <Space size="small" wrap>
                    <Tooltip title="粗体">
                      <Button
                        size="small"
                        icon={<BoldOutlined />}
                        onClick={() => handleFormatClick(() => editor.commands.toggleBold())}
                        type={editor.isActive('bold') ? 'primary' : 'default'}
                      />
                    </Tooltip>

                    <Tooltip title="斜体">
                      <Button
                        size="small"
                        icon={<ItalicOutlined />}
                        onClick={() => handleFormatClick(() => editor.commands.toggleItalic())}
                        type={editor.isActive('italic') ? 'primary' : 'default'}
                      />
                    </Tooltip>

                    <Tooltip title="删除线">
                      <Button
                        size="small"
                        icon={<StrikethroughOutlined />}
                        onClick={() => handleFormatClick(() => editor.commands.toggleStrike())}
                        type={editor.isActive('strike') ? 'primary' : 'default'}
                      />
                    </Tooltip>

                    <div style={{ width: '1px', height: '24px', background: '#d9d9d9' }} />

                    <Tooltip title="无序列表">
                      <Button
                        size="small"
                        icon={<UnorderedListOutlined />}
                        onClick={() => handleFormatClick(() => editor.commands.toggleBulletList())}
                        type={editor.isActive('bulletList') ? 'primary' : 'default'}
                      />
                    </Tooltip>

                    <Tooltip title="有序列表">
                      <Button
                        size="small"
                        icon={<OrderedListOutlined />}
                        onClick={() => handleFormatClick(() => editor.commands.toggleOrderedList())}
                        type={editor.isActive('orderedList') ? 'primary' : 'default'}
                      />
                    </Tooltip>

                    <div style={{ width: '1px', height: '24px', background: '#d9d9d9' }} />

                    <Tooltip title="标题 1">
                      <Button
                        size="small"
                        onClick={() => handleFormatClick(() => editor.commands.toggleHeading({ level: 1 }))}
                        type={editor.isActive('heading', { level: 1 }) ? 'primary' : 'default'}
                      >
                        H1
                      </Button>
                    </Tooltip>

                    <Tooltip title="标题 2">
                      <Button
                        size="small"
                        onClick={() => handleFormatClick(() => editor.commands.toggleHeading({ level: 2 }))}
                        type={editor.isActive('heading', { level: 2 }) ? 'primary' : 'default'}
                      >
                        H2
                      </Button>
                    </Tooltip>

                    <Tooltip title="标题 3">
                      <Button
                        size="small"
                        onClick={() => handleFormatClick(() => editor.commands.toggleHeading({ level: 3 }))}
                        type={editor.isActive('heading', { level: 3 }) ? 'primary' : 'default'}
                      >
                        H3
                      </Button>
                    </Tooltip>

                    <div style={{ width: '1px', height: '24px', background: '#d9d9d9' }} />

                    <Tooltip title="撤销">
                      <Button
                        size="small"
                        icon={<UndoOutlined />}
                        onClick={() => editor.commands.undo()}
                        disabled={!editor.can().undo()}
                      />
                    </Tooltip>

                    <Tooltip title="重做">
                      <Button
                        size="small"
                        icon={<RedoOutlined />}
                        onClick={() => editor.commands.redo()}
                        disabled={!editor.can().redo()}
                      />
                    </Tooltip>

                    <div style={{ width: '1px', height: '24px', background: '#d9d9d9' }} />

                    <Tooltip title="清空">
                      <Button size="small" icon={<ClearOutlined />} onClick={clearContent} />
                    </Tooltip>
                  </Space>
                </div>

                {/* 编辑器内容区 */}
                <div className="tiptap-content" style={{ height }}>
                  <EditorContent editor={editor} />
                </div>
              </div>
            ),
          },
          {
            key: 'markdown',
            label: 'Markdown 源码',
            children: (
              <div>
                <div style={{ marginBottom: 12 }}>
                  <Space>
                    <Button size="small" icon={<CopyOutlined />} onClick={copyMarkdown}>
                      复制内容
                    </Button>
                  </Space>
                </div>
                <textarea
                  value={value}
                  onChange={(e) => {
                    onChange(e.target.value)
                    editor.commands.setContent(e.target.value)
                  }}
                  style={{
                    width: '100%',
                    height,
                    padding: '12px',
                    border: '1px solid #d9d9d9',
                    borderRadius: '4px',
                    fontFamily: 'monospace',
                    fontSize: '12px',
                  }}
                  placeholder={placeholder}
                />
              </div>
            ),
          },
          {
            key: 'preview',
            label: '预览',
            children: (
              <div
                style={{
                  height,
                  padding: '12px',
                  border: '1px solid #d9d9d9',
                  borderRadius: '4px',
                  backgroundColor: '#fafafa',
                  overflowY: 'auto',
                }}
              >
                <div
                  dangerouslySetInnerHTML={{ __html: value }}
                  style={{
                    fontSize: '14px',
                    lineHeight: '1.8',
                  }}
                />
              </div>
            ),
          },
        ]}
      />
    </div>
  )
}
