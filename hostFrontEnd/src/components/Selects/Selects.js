import React, { PureComponent } from 'react'
import { Select } from 'antd'
import PropTypes from 'prop-types'

class Selects extends PureComponent {
  handleChange = value => {
    this.props.onChange(value)
  }

  render() {
    const children = []
    for (let i = 0; i < this.props.totalIds; i++) {
      children.push(
        <Select.Option key={this.props.totalIds[i]}>
          {this.props.totalIds[i]}
        </Select.Option>
      )
    }

    return (
      <Select
        mode="tags"
        style={{ width: '100%' }}
        placeholder="选择"
        defaultValue={this.props.ids}
        onChange={this.handleChange}
        disabled={this.props.disabled}
      >
        {children}
      </Select>
    )
  }
}

Selects.propTypes = {
  ids: PropTypes.array,
  totalIds: PropTypes.array,
}

export default Selects
