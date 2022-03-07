import axios from 'axios'
import { useRouter } from 'next/router'

const useSearchDocument = ({ setHits, setNotFound, setIsLoading }) => {
  const router = useRouter()
  const fetchData = (val) => {

    if (val == "") {
      return
    }

    // if (filter == null) {

    //   if (router.query["filter[repo]"]) {
    //     const newFilter = Object.assign(filter || {}, {
    //       repo: router.query["filter[repo]"].split(",")
    //     });
    //     filter = newFilter;
    //   }

    //   if (router.query["filter[lang]"]) {
    //     const newFilter = Object.assign(filter || {}, {
    //       lang: router.query["filter[lang]"].split(",")
    //     });
    //     filter = newFilter;
    //   }

    //   if (router.query["filter[path]"]) {
    //     const newFilter = Object.assign(filter || {}, {
    //       path: router.query["filter[path]"].split(",")
    //     });
    //     filter = newFilter;
    //   }

    //   if (filter != null) {
    //     setFilter(() => filter);
    //   }

    // }

    let queryParam = val
    // if (filter?.repo?.length > 0) {
    //   queryParam += `&filter[repo]=${filter.repo.join(",")}`
    // }

    // if (filter?.lang?.length > 0) {
    //   queryParam += `&filter[lang]=${filter.lang.join(",")}`
    // }

    // if (filter?.path?.length > 0) {
    //   queryParam += `&filter[path]=${filter.path.join(",")}`
    // }

    setIsLoading(true)
    // axios.get(`https://heline.dev/api/search?q=${queryParam}&tbm=docs`)
    console.log(`/api/search?q=${queryParam}&tbm=docs`)
    axios.get(`http://localhost:8000/api/search?q=${queryParam}&tbm=docs`)
      .then(res => {
        setIsLoading(false)
        if (res.data.docs.hits === null) {
          setNotFound(true)
          setHits(null)
        } else {
          setNotFound(false)
          setHits(res.data.docs)
        }
      })
      .catch(e => {
        setIsLoading(false)
        setNotFound(true)
        setHits(null)
        console.log(e)
      })

    if (typeof window !== 'undefined') {
      queryParam = queryParam + "&tbm=docs"
      router.push(`search?q=${queryParam}`)
    }
  }

  return fetchData
}

export default useSearchDocument