import Zoom from 'react-medium-image-zoom'
import 'react-medium-image-zoom/dist/styles.css'
import './zoom-custom.css';

export const ZoomableImage = ({ src, alt }: { src: string; alt?: string }) => {
    return (
        <Zoom>
            <img
                src={src}
                alt={alt}
                style={{ width: '300px', borderRadius: '8px', cursor: 'pointer' }}
            />
        </Zoom>
    )
}