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
class ClassromModal extends PureComponent {
  handleOk = () => {
    const { onOk, form } = this.props
    const { validateFields, getFieldsValue } = form

    validateFields(errors => {
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
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
          <FormItem label={`教室号`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('classroom_no', {
              initialValue: item.classroom_no,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input disabled={ type === 'create' ? false : true } />)}
          </FormItem>
          <FormItem label={`摄像头ID`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('camera_id', {
              initialValue: item.camera_id,
              rules: [
              ],
            })(<Input />)}
          </FormItem>
        </Form>
      </Modal>
    )
  }
}

ClassromModal.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default ClassromModal
