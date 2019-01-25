import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Modal } from 'antd'
import { Selects, PicturesWall } from 'components'

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
class ClassModal extends PureComponent {
  handleOk = () => {
    const { item = {}, onOk, form } = this.props
    const { validateFields, getFieldsValue } = form

    validateFields(errors => {
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
        class_id: item.class_id,
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
          <FormItem label={`班级名`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('class_name', {
              initialValue: item.class_name,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`教室号`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('classroom_no', {
              initialValue: item.classroom_no,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`教师`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('teacher_nos', {
              initialValue: item.teacher_nos,
              valuePropName: "ids",
            })(<Selects />)}
          </FormItem>
          <FormItem label={`学生`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('student_nos', {
              initialValue: item.student_nos,
              valuePropName: "ids",
            })(<Selects disabled={ type === 'create' ? true : false } />)}
          </FormItem>
          <FormItem label={`班级图`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('class_image', {
              initialValue: item.class_image,
              valuePropName: "image",
            })(<PicturesWall />)}
          </FormItem>
        </Form>
      </Modal>
    )
  }
}

ClassModal.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default ClassModal
