package utils

// Reflection -> Reflect
// Create map[Key_Type]Value_Type{} to use structs as map keys -> Useful to learn deeper structures
// Work with trees, linkedLists & more based on Structs

// Method that receives a Struct (coming from PUT, i.e. it can be partial) & the old Struct (living in DB) & returns
// a final Struct to be saved with only the right fields updated -> Avoiding null insertion
