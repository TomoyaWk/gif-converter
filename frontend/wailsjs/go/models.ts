export namespace main {
	
	export class GifConvertOptions {
	    FPS: number;
	    Width: number;
	    Height: number;
	    Quality: string;
	
	    static createFrom(source: any = {}) {
	        return new GifConvertOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FPS = source["FPS"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	        this.Quality = source["Quality"];
	    }
	}

}

