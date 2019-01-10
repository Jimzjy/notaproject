import { PureComponent } from "react";
import { Select } from "antd";
import PropTypes from 'prop-types'

class Selects extends PureComponent {
  render() {
    const children = [];
    for (let i = 0; i < this.props.totalClasses; i++) {
      children.push(<Option key={this.props.totalClasses[i]}>{this.props.totalClasses[i]}</Option>);
    }

    return (
      <Select
        mode="multiple"
        style={{ width: '100%' }}
        placeholder="选择班级"
        defaultValue={this.props.classIds}
        onChange={this.props.onChange}
      >
        {children}
      </Select>
    )
  }
}

Selects.propTypes = {
  classIds: PropTypes.array,
  totalClasses: PropTypes.array,
}

export default Selects