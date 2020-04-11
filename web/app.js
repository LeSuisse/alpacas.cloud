import SwaggerUI from "swagger-ui";
import "swagger-ui/dist/swagger-ui.css";

document.addEventListener("DOMContentLoaded", () => {
    SwaggerUI({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        validatorUrl: null,
        docExpansion: 'full',
        layout: "BaseLayout",
    });
});
