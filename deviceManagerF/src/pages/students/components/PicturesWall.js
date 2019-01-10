import { Upload, Icon, Modal } from 'antd';
import { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { apiPrefix } from 'utils/config'

class PicturesWall extends PureComponent {
  state = {
    previewVisible: false,
    previewImage: '',
    fileList: [],
  };

  componentDidMount() {
    let list = this.props.image === undefined ? [] : [{
      uid: '-1',
      name: 'xxx.png',
      status: 'done',
      url: `${apiPrefix}/images/${this.props.image}`,
    }]
    
    this.setState({
      fileList: list
    })
  }

  handleCancel = () => this.setState({ previewVisible: false })

  handlePreview = (file) => {
    this.setState({
      previewImage: file.url,
      previewVisible: true,
    });
  }

  handleChange = ({ fileList, file }) => {
    if (file.status === "done") {
      let imageName = fileList.length > 0 ? fileList[0].response : {message: ''}
      this.props.onChange(imageName.message)
    }
    
    this.setState({ fileList })
  }

  render() {
    const { previewVisible, previewImage, fileList } = this.state;
    const uploadButton = (
      <div>
        <Icon type="plus" />
        <div className="ant-upload-text">上传</div>
      </div>
    );
    return (
      <div className="clearfix">
        <Upload
          action={`${apiPrefix}/images`}
          listType="picture-card"
          fileList={fileList}
          onPreview={this.handlePreview}
          onChange={this.handleChange}
        >
          {fileList.length >= 1 ? null : uploadButton}
        </Upload>
        <Modal visible={previewVisible} footer={null} onCancel={this.handleCancel}>
          <img style={{ width: '100%' }} src={previewImage} />
        </Modal>
      </div>
    );
  }
}

PicturesWall.propTypes = {
  image: PropTypes.string
}

export default PicturesWall