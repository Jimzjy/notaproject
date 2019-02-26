import React, { PureComponent } from 'react'
import { connect } from 'dva'
import { Row, Col, Card } from 'antd'
import { Page } from 'components'
import { apiPrefix } from 'utils/config'
import { router } from 'utils'

@connect(({ standup, loading }) => ({
  standup,
  loading,
}))
class StandUp extends PureComponent {
  handleClassCardClick = id => {
    router.push(`/standup/${id}`)
  }

  render() {
    const { standup, loading } = this.props
    const { classes } = standup

    const classCards = classes.map((item, key) => (
      <Col key={key} lg={6} md={6}>
        <Card
          hoverable
          cover={<img src={`${apiPrefix}/images/${item.class_image}`} />}
          onClick={() => {
            this.handleClassCardClick(item.class_id)
          }}
        >
          <Card.Meta
            title={item.class_name}
            description={`ID: ${item.class_id}`}
          />
        </Card>
      </Col>
    ))

    return (
      <Page>
        <Row gutter={24}>{classCards}</Row>
      </Page>
    )
  }
}

export default StandUp
