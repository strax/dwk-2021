package org.strax.dwk;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.MediaType;
import java.io.File;

@Path("/image")
public class ImageResource {
    @GET
    @Produces("image/jpg")
    public Response get() {
        var file = new File("/mnt/volume1/image.jpeg");
        return Response.ok(file).build();
    }
}
