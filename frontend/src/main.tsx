import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import Videos from './components/videos.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Videos />
  </StrictMode>,
)
