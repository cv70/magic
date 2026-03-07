import React from 'react'
import { Card, Row, Col, Tooltip } from 'antd'
import type { Domain } from '../constants/domains'
import { getDomainsList } from '../constants/domains'
import styles from './DomainSelector.module.css'

interface DomainSelectorProps {
  value?: string
  onChange?: (domainId: string) => void
  size?: 'small' | 'large'
  showDescription?: boolean
}

/**
 * 领域选择组件
 * 展示所有可用领域，支持点击选择
 */
export const DomainSelector: React.FC<DomainSelectorProps> = ({
  value,
  onChange,
  size = 'large',
  showDescription = true
}) => {
  const domains = getDomainsList()

  const handleDomainClick = (domainId: string) => {
    onChange?.(domainId)
  }

  const isSmall = size === 'small'

  return (
    <div className={styles.container}>
      <Row gutter={[12, 12]}>
        {domains.map((domain: Domain) => (
          <Col
            key={domain.id}
            xs={12}
            sm={isSmall ? 12 : 8}
            md={isSmall ? 8 : 6}
            lg={isSmall ? 6 : 4}
          >
            <Tooltip title={domain.description} placement="top">
              <Card
                hoverable
                onClick={() => handleDomainClick(domain.id)}
                className={`${styles.domainCard} ${
                  value === domain.id ? styles.selected : ''
                }`}
                style={{
                  borderColor: value === domain.id ? domain.color : 'transparent',
                  borderWidth: '2px'
                }}
              >
                <div className={styles.icon}>{domain.icon}</div>
                <div className={styles.name}>{domain.name}</div>
                {showDescription && (
                  <div className={styles.description}>{domain.description}</div>
                )}
              </Card>
            </Tooltip>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default DomainSelector
