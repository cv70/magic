import React, { useState } from 'react'
import { Form, Input, InputNumber, Select, Slider, Button, Space } from 'antd'
import type { DomainTemplate, ParamField } from '../constants/domainTemplates'

interface DomainParamsFormProps {
  template: DomainTemplate
  onSubmit?: (params: Record<string, any>) => void
  loading?: boolean
}

/**
 * 领域参数表单组件
 * 根据模板配置动态渲染参数输入字段
 */
export const DomainParamsForm: React.FC<DomainParamsFormProps> = ({
  template,
  onSubmit,
  loading = false
}) => {
  const [form] = Form.useForm()
  const [params, setParams] = useState<Record<string, any>>({})

  const handleParamChange = (key: string, value: any) => {
    setParams(prev => ({ ...prev, [key]: value }))
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      onSubmit?.(values)
    } catch (error) {
      console.error('Form validation failed:', error)
    }
  }

  const renderField = (field: ParamField) => {
    const key = field.key
    const commonProps = {
      value: params[key] !== undefined ? params[key] : field.defaultValue,
      onChange: (value: any) => {
        handleParamChange(key, value)
      }
    }

    switch (field.type) {
      case 'text':
        return (
          <Input
            placeholder={field.placeholder}
            {...commonProps}
          />
        )
      case 'textarea':
        return (
          <Input.TextArea
            placeholder={field.placeholder}
            rows={4}
            {...commonProps}
          />
        )
      case 'select':
        return (
          <Select
            placeholder={field.placeholder}
            options={field.options}
            {...commonProps}
          />
        )
      case 'number':
        return (
          <InputNumber
            min={field.min}
            max={field.max}
            placeholder={field.placeholder}
            style={{ width: '100%' }}
            {...commonProps}
          />
        )
      case 'slider':
        return (
          <div>
            <Slider
              min={field.min || 0}
              max={field.max || 100}
              {...commonProps}
            />
            <div style={{ marginTop: '8px', textAlign: 'center', fontSize: '12px', color: '#999' }}>
              当前: {params[key] !== undefined ? params[key] : field.defaultValue}
            </div>
          </div>
        )
      default:
        return null
    }
  }

  return (
    <div>
      <Form
        form={form}
        layout="vertical"
        autoComplete="off"
      >
        {template.paramFields.map((field: ParamField) => (
          <Form.Item
            key={field.key}
            label={field.label}
            name={field.key}
            required={field.required}
            rules={
              field.required
                ? [{ required: true, message: `请填写${field.label}` }]
                : []
            }
            initialValue={field.defaultValue}
          >
            {renderField(field)}
          </Form.Item>
        ))}
      </Form>

      <Space style={{ marginTop: '24px', width: '100%', justifyContent: 'flex-end' }}>
        <Button loading={loading} type="primary" size="large" onClick={handleSubmit}>
          生成内容
        </Button>
      </Space>
    </div>
  )
}

export default DomainParamsForm
