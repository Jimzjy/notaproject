import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Modal } from 'antd'
import { Selects } from 'components'

const FormItem = Form.Item

const formItemLayout = {
  labelCol: {
    span: 6,
  },
  wrapperCol: {
    span: 14,
  },
}

@Form.create()
class DeviceModal extends PureComponent {
  handleOk = () => {
    const { item = {}, onOk, form } = this.props
    const { validateFields, getFieldsValue } = form

    validateFields(errors => {
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
        device_id: item.device_id,
      }
      onOk(data)
    })
  }

  render() {
    const { item = {}, onOk, form, type, ...modalProps } = this.props
    const { getFieldDecorator } = form

    return (
      <Modal {...modalProps} onOk={this.handleOk}>
        <Form layout="horizontal">
          <FormItem label={`设备IP`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('device_path', {
              initialValue: item.device_path,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`设备Port`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('device_port', {
              initialValue: item.device_port,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`摄像头`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('camera_ids', {
              initialValue: item.camera_ids,
              valuePropName: "ids",
            })(<Selects />)}
          </FormItem>
        </Form>
      </Modal>
    )
  }
}

DeviceModal.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default DeviceModal
