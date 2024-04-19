vcl 4.1;
import std;

backend default {
    .host = "envoy-service";
    .port = "9096";
}

acl purge {
    "localhost";
    "envoy-service";
}

sub vcl_recv {
     # after adding ext_authz filter
     unset req.http.authorization;

     set req.http.Host = req.http.redirect-backend;
     set req.url = req.http.x-temp-path;
     set req.url = std.querysort(req.url);


     if (req.method == "PURGE") {
		# check if the client is allowed to purge content
		if (!client.ip ~ purge) {
			return(synth(405,"Not allowed."));
		}
		return (purge);
	}
}


sub vcl_miss {
    if (req.http.x-cluster-header == "varnish_backend_cluster") {
        set req.http.x-cluster-header = req.http.redirect-backend;  
    } 
        return (fetch);
    
}
#set the header for requests by passing the cache, ex: POST,PUT
sub vcl_pass{
    if (req.http.x-cluster-header == "varnish_backend_cluster") {
        set req.http.x-cluster-header = req.http.redirect-backend;  
    } 
}

sub vcl_backend_fetch {

    if (bereq.http.x-temp-path) {
        set bereq.url = bereq.http.x-temp-path;
        set bereq.url = std.querysort(bereq.url);
    }
}

sub vcl_backend_response {
    # Don't cache 404 responses
    if (beresp.status == 404) {
        set beresp.uncacheable = true;
    }

    # Determine the appropriate storage and check response size against available storage size
    if (bereq.url ~ "/org1/") {
        if (std.integer(beresp.http.Content-Length, 0) > std.integer(std.getenv("ORG1_CACHE_MAX_SIZE"), 0)) {
            set beresp.uncacheable = true;
            return (deliver);
        }
        # set beresp.storage = storage.org1;
        # set beresp.http.x-storage = "org1";
    } elsif (bereq.url ~ "/org2/") {
        if (std.integer(beresp.http.Content-Length, 0) > std.integer(std.getenv("ORG2_CACHE_MAX_SIZE"), 0)) {
            set beresp.uncacheable = true;
            return (deliver);
        }
        # set beresp.storage = storage.org2;
        # set beresp.http.x-storage = "org2";
    } else {
          if (std.integer(beresp.http.content-length, 0) > 1048576) {
            set beresp.uncacheable = true;
            return (deliver);
        }
        set beresp.storage = storage.default; // Default storage
        set beresp.http.x-storage = "default";
    }
}




sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cached-By = "Varnish";
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }
}



