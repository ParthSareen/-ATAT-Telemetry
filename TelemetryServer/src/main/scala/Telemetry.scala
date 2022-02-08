import akka.actor.ActorSystem
import akka.actor.Status.{Failure, Success}
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.StatusCodes
import akka.http.scaladsl.server.Directives
import akka.http.scaladsl.server.Directives._
import akka.stream.ActorMaterializer

object Telemetry extends App {
  // TODO: Move to config
  val host = "localhost"
  val port = 7000

  implicit val system = ActorSystem("HTTP_SERVER")
  implicit val mat : ActorMaterializer.type = ActorMaterializer
    import system.dispatcher

  val route =
      Directives.get{
        path("hello"){
          complete(StatusCodes.OK, "Hello from Server")
        }
    }

  val binding = Http().newServerAt(host, port).bindFlow(route)
  println(s"Server now online. Please navigate to http://localhost:8080/hello\nPress RETURN to stop...")
  binding
    .flatMap(_.unbind())
    .onComplete(_ => system.terminate())

}
