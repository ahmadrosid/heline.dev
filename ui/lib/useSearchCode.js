import axios from 'axios'
import { useRouter } from 'next/router'

const useSearchCode = ({ setHits, setNotFound, setIsLoading, setFilter }) => {
  const router = useRouter()
  const fetchData = (val, filter) => {

    if (val == "") {
      return
    }

    if (filter == null) {

      if (router.query["filter[repo]"]) {
        const newFilter = Object.assign(filter || {}, {
          repo: router.query["filter[repo]"].split(",")
        });
        filter = newFilter;
      }

      if (router.query["filter[lang]"]) {
        const newFilter = Object.assign(filter || {}, {
          lang: router.query["filter[lang]"].split(",")
        });
        filter = newFilter;
      }

      if (router.query["filter[path]"]) {
        const newFilter = Object.assign(filter || {}, {
          path: router.query["filter[path]"].split(",")
        });
        filter = newFilter;
      }

      if (filter != null) {
        setFilter(() => filter);
      }

    }

    let queryParam = val
    if (filter?.repo?.length > 0) {
      queryParam += `&filter[repo]=${filter.repo.join(",")}`
    }

    if (filter?.lang?.length > 0) {
      queryParam += `&filter[lang]=${filter.lang.join(",")}`
    }

    if (filter?.path?.length > 0) {
      queryParam += `&filter[path]=${filter.path.join(",")}`
    }

    setIsLoading(true)
    // Use environment variable for API URL
      axios.get(`${process.env.NEXT_PUBLIC_API_URL}/api/search?q=${queryParam}`)
      .then(res => {
        setIsLoading(false)
        if (res.data.hits.hits === null) {
          setNotFound(true)
          setHits(null)
        } else {
          setNotFound(false)
          setHits(res.data.hits)
        }
      })
      .catch(e => {
        setIsLoading(false)
        setNotFound(true)
        setHits(null)
      })

    if (typeof window !== 'undefined') {
      router.push(`search?q=${queryParam}`)
    }
  }

  return fetchData
}

export default useSearchCode