import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Modal } from 'antd'

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
class CameraModal extends PureComponent {
  handleOk = () => {
    const { item = {}, onOk, form } = this.props
    const { validateFields, getFieldsValue } = form

    validateFields(errors => {
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
        camera_id: item.camera_id,
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
          <FormItem label={`Stream 地址`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('cam_stream_path', {
              initialValue: item.cam_stream_path,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`ONVIF 地址`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('cam_onvif_path', {
              initialValue: item.cam_onvif_path,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`ONVIF 用户名`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('cam_auth_name', {
              initialValue: item.cam_auth_name,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`ONVIF 密码`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('cam_auth_password', {
              initialValue: item.cam_auth_password,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`设备ID`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('device_id', {
              initialValue: item.device_id,
              rules: [
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`教室号`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('classroom_no', {
              initialValue: item.classroom_no,
              rules: [
              ],
            })(<Input />)}
          </FormItem>
        </Form>
      </Modal>
    )
  }
}

CameraModal.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default CameraModal
