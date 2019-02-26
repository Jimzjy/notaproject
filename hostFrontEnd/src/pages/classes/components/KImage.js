import React from 'react'
import { Image } from 'react-konva'
import Konva from 'konva'

class KImage extends React.Component {
  state = {
    image: null,
  }
  componentDidMount() {
    const image = new window.Image()
    image.src = this.props.src
    image.onload = () => {
      this.setState({
        image: image,
      })
    }
  }

  render() {
    return (
      <Image
        image={this.state.image}
        width={this.props.width}
        height={this.props.height}
      />
    )
  }
}

export default KImage
