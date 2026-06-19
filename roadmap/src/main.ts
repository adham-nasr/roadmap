import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { NestFactory } from '@nestjs/core';
import { ExpressAdapter } from '@nestjs/platform-express';
import { AppModule } from './app.module';

import express from 'express';
import serverlessExpress from '@codegenie/serverless-express';

let server: any;

async function bootstrap() {
    console.log('BOOTSTRAP START');
  const expressApp = express();

  const app = await NestFactory.create(
    AppModule,
    new ExpressAdapter(expressApp),
  );

    console.log('NEST CREATED');
  const options = new DocumentBuilder()
    .setTitle('Roadmap')
    .setDescription('The Roadmap API description')
    .setVersion('1.0')
    .addTag('Roadmaps')
    .build();

  const document = SwaggerModule.createDocument(app, options);
  SwaggerModule.setup('api', app, document);

  await app.init();

    console.log('BOOTSTRAP DONE');
  return serverlessExpress({
    app: expressApp,
  });
}

export const handler = async (event: any, context: any) => {
  if (!server) {
    server = await bootstrap();
  }

  return await server(event, context);
};