import axios from 'axios'
import { useRouter } from 'next/router'

const useGetDocumentByID = ({ setHits, setNotFound, setIsLoading }) => {
  const router = useRouter()

  const fetchData = (id) => {

    if (id == "") {
      return
    }

    let queryParam = `?tbm=docs&id=${id}`
    setIsLoading(true)

    // axios.get(`https://heline.dev/api/search?q=${queryParam}&tbm=docs`)
    axios.get(`http://localhost:8000/api/search${queryParam}`)
      .then(res => {
        setIsLoading(false)
        if (res.data === null) {
          setNotFound(true)
          setHits(null)
        } else {
          setNotFound(false)
          setHits(res.data)
        }
      })
      .catch(e => {
        setIsLoading(false)
        setNotFound(true)
        setHits(null)
        console.log(e)
      })

    if (typeof window !== 'undefined') {
      // router.push(`search${queryParam}`)
    }
  }

  return fetchData
}

export default useGetDocumentByID