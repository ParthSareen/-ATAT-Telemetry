ThisBuild / version := "0.1.0-SNAPSHOT"

ThisBuild / scalaVersion := "2.13.8"
lazy val root = (project in file("."))
  .settings(
    name := "TelemetryServer"
  )

Compile / PB.targets := Seq(
  scalapb.gen() -> (Compile / sourceManaged).value / "scalapb"
)

val AkkaVersion = "2.6.18"
val AkkaHttpVersion = "10.2.7"
val akka = Seq(
  "com.typesafe.akka" %% "akka-actor-typed" % AkkaVersion,
  "com.typesafe.akka" %% "akka-stream" % AkkaVersion,
  "com.typesafe.akka" %% "akka-http" % AkkaHttpVersion
)
libraryDependencies ++= Seq() ++ akka
