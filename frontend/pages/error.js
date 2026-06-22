// <head>
//     <meta charset="UTF-8">
//     <meta name="viewport" content="width=device-width, initial-scale=1.0">
//     <title>Error {{.Code}}</title>
//     <link rel="stylesheet" href="/assets/style.css"> 
// </head>

import { Error } from '../components/Error.js';

export async function render() {
    return Error(200, "hello"); // just example !
}
