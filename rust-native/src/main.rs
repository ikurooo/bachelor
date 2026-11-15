use axum::{
    extract::Request,
    http::{header, StatusCode},
    response::IntoResponse,
    routing::any,
    Router,
};

#[tokio::main]
async fn main() {
    let app = Router::new().fallback(any(handle_hello_rust));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .unwrap();
    
    println!("Listening on http://{}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}

async fn handle_hello_rust(req: Request) -> impl IntoResponse {
    println!("Handling request to {:?}", req.uri());

    (
        StatusCode::OK,
        [(header::CONTENT_TYPE, "text/plain")],
        "Hello World!",
    )
}