import { createClient } from '@supabase/supabase-js'

async function fetchProducts() {
    const response = await fetch("http://localhost/backend/main.go")   
    const data = await response.json()

}